package k8s

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/glory-go/glory/common"
	"github.com/glory-go/glory/config"
	"github.com/glory-go/glory/log"
	"github.com/glory-go/glory/plugin"
	"github.com/glory-go/glory/registry"
	"github.com/glory-go/glory/tools"
	perrors "github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

const (
	// defined in yaml
	podNameKey   = "HOSTNAME"
	nameSpaceKey = "NAMESPACE"

	// glory labelKey
	GloryLabelKey   = "glory"
	GloryLabelValue = "online"
)

type k8sRegistry struct {
	client          *kubernetes.Clientset
	podName         string
	podNamespace    string
	eventHandlerMap map[string]*eventHandler
}

func init() {
	plugin.SetRegistryFactory("k8s", newK8sRegistry)
}

func newK8sRegistry(registryConfig *config.RegistryConfig) registry.Registry {
	clusterCfg, err := rest.InClusterConfig()
	if err != nil {
		log.Errorf("K8s registry In Cluster config get error = %v", err)
		return nil
	}

	// 根据指定的 config 创建一个新的 clientset
	clientset, err := kubernetes.NewForConfig(clusterCfg)
	if err != nil {
		panic(err.Error())
	}

	podname := os.Getenv(podNameKey)
	podnamespace := os.Getenv(nameSpaceKey)

	return &k8sRegistry{
		client:          clientset,
		podName:         podname,
		podNamespace:    podnamespace,
		eventHandlerMap: make(map[string]*eventHandler),
	}
}

func (kr *k8sRegistry) Subscribe(key string) (chan common.RegistryChangeEvent, error) {
	log.Debugf("in subscribe with key = %s, namespace = %s, name = %s\n", key, kr.podNamespace, kr.podName)
	req, err := labels.NewRequirement(GloryLabelKey, selection.In, []string{GloryLabelValue})
	if err != nil {
		return nil, perrors.WithMessage(err, "new requirement")
	}

	informersFactory := informers.NewSharedInformerFactoryWithOptions(
		kr.client,
		5*time.Minute, // todo configurable
		informers.WithNamespace(kr.podNamespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.LabelSelector = req.String()
		}),
	)
	stopper := make(chan struct{})
	go informersFactory.Start(stopper)
	podInformer := informersFactory.Core().V1().Pods()

	eventChan := make(chan common.RegistryChangeEvent)

	handler := newK8sEventHandler(eventChan, key, stopper)
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    handler.add,
		UpdateFunc: handler.update,
		DeleteFunc: handler.delete,
	})
	if !cache.WaitForCacheSync(context.Background().Done(), podInformer.Informer().HasSynced) {
		log.Errorf("wait for cache sync finish @namespace %s fail", kr.podNamespace)
		return nil, nil
	}
	kr.eventHandlerMap[key] = handler
	log.Debug("subscribe success!")
	return eventChan, nil
}

func (kr *k8sRegistry) Unsubscribe(key string) error {
	close(kr.eventHandlerMap[key].eventChan)
	kr.eventHandlerMap[key].stopper <- struct{}{}
	delete(kr.eventHandlerMap, key)
	return nil
}

func (kr *k8sRegistry) Register(serviceID string, localAddress common.Address) {
	log.Debugf("register service id = %s\n", serviceID)

	oldPod := &v1.Pod{}
	newPod := &v1.Pod{}
	oldPod.Labels = make(map[string]string, 8)
	newPod.Labels = make(map[string]string, 8)

	newPod.Labels[serviceID] = tools.Addr2AddrLabel(localAddress)
	newPod.Labels[GloryLabelKey] = GloryLabelValue

	oldData, err := json.Marshal(oldPod)
	if err != nil {
		panic(err)
	}

	newData, err := json.Marshal(newPod)
	if err != nil {
		panic(err)
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Pod{})
	if err != nil {
		panic(err)
	}

	_, err = kr.client.CoreV1().Pods(kr.podNamespace).Patch(context.Background(), kr.podName, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}

func (nr *k8sRegistry) UnRegister(serviceID string, localAddress common.Address) {
	cfg, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	c, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err)
	}

	oldPod := &v1.Pod{}
	newPod := &v1.Pod{}
	oldPod.Labels = make(map[string]string, 8)
	newPod.Labels = make(map[string]string, 8)

	oldPod.Labels[serviceID] = tools.Addr2AddrLabel(localAddress)
	oldPod.Labels[GloryLabelKey] = GloryLabelValue
	oldData, err := json.Marshal(oldPod)
	if err != nil {
		panic(err)
	}

	newData, err := json.Marshal(newPod)
	if err != nil {
		panic(err)
	}

	patchBytes, err := strategicpatch.CreateTwoWayMergePatch(oldData, newData, v1.Pod{})
	if err != nil {
		panic(err)
	}

	_, err = c.CoreV1().Pods(nr.podNamespace).Patch(context.Background(), nr.podName, types.StrategicMergePatchType, patchBytes, metav1.PatchOptions{})
	if err != nil {
		panic(err)
	}
}

func (nr *k8sRegistry) Refer(key string) []common.Address {
	req, err := labels.NewRequirement(key, selection.Exists, []string{})
	if err != nil {
		log.Error("k8s registry refer err = ", err)
		return []common.Address{}
	}
	pods, err := nr.client.CoreV1().Pods(nr.podNamespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: req.String(),
	})
	if err != nil {
		return []common.Address{}
	}
	result := make([]common.Address, 0)
	for _, v := range pods.Items {
		result = append(result, tools.AddrLabel2Addr(v.Labels[key]))
	}
	log.Debug("refer provider list for ", key, " is ", result)
	return result
}
