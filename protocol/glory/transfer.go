package glory

import (
	"encoding/binary"
	"errors"
	"io"
	"net"
	"syscall"
)

import (
	"github.com/glory-go/glory/log"
)

const (
	defaultMaxBufferSize = 16384
	//defualtMaxQueueSize  = 1024
)

type GloryConnClientStatus uint32

const (
	GloryConnClientHealthy = GloryConnClientStatus(0)
	//GloryConnClientClosed  = GloryConnClientStatus(1)
)

type gloryConnClient struct {
	conn          *net.TCPConn
	maxBufferSize uint32
	frameQueue    chan []byte
	status        GloryConnClientStatus
}

func newGloryConnClient(targetAddress string) (*gloryConnClient, error) {
	addr, err := net.ResolveTCPAddr("tcp", targetAddress)
	if err != nil {
		log.Error("new glory invoker:net.ResolveTCPAddr failed with err = ", err, " address = ", targetAddress)
		return nil, err
	}
	log.Debugf("glory conn dail addr = ", targetAddress)
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Error("new glory invoker:net.DialTCP failed with err = ", err, " address = ", addr)
		return nil, err
	}
	newClient := &gloryConnClient{
		conn:          conn,
		maxBufferSize: defaultMaxBufferSize,
		frameQueue:    make(chan []byte, defaultMaxBufferSize),
		status:        GloryConnClientHealthy,
	}
	go newClient.runReceivingData()
	return newClient, nil
}

func newGloryConnClientFromConn(conn *net.TCPConn) *gloryConnClient {
	newClient := &gloryConnClient{
		conn:          conn,
		maxBufferSize: defaultMaxBufferSize,
		frameQueue:    make(chan []byte, defaultMaxBufferSize),
		status:        GloryConnClientHealthy,
	}
	go newClient.runReceivingData()
	return newClient
}

// ReadFrame can split data frame by glory header
// length | data
// 4 byte | data
func (cc *gloryConnClient) ReadFrame() ([]byte, error) {
	data := <-cc.frameQueue
	if data == nil {
		return nil, errors.New("read on closed client")
	}
	return data, nil
}

func (cc *gloryConnClient) WriteFrame(data []byte) (int, error) {
	return cc.conn.Write(data2Frame(data))
}

func frame2Datas(fm []byte) [][]byte {
	result := make([][]byte, 0)
	for len(fm) > 0 {
		h := fm[:4]
		length := binary.BigEndian.Uint32(h)
		result = append(result, fm[4:4+length])
		fm = fm[4+length:]
	}
	return result
}

func data2Frame(data []byte) []byte {
	length := uint32(len(data))
	fm := make([]byte, 4+length)
	binary.BigEndian.PutUint32(fm, length)
	copy(fm[4:], data)
	return fm
}

func (cc *gloryConnClient) runReceivingData() {
	for {
		buf := make([]byte, defaultMaxBufferSize)
		n, err := cc.conn.Read(buf)
		if err != nil {
			if err == syscall.EINVAL || err == io.EOF {
				// conn closed
				close(cc.frameQueue)
				return
			}
			log.Error("gloryConnClient read from conn err = ", err)
			continue
		}
		buf = buf[0:n]
		for _, v := range frame2Datas(buf) {
			cc.frameQueue <- v
		}
	}
}

func (cc *gloryConnClient) Close() {
	cc.conn.Close()
}
