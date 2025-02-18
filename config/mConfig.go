package config

import (
	"aDi/log"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"reflect"
	"sync/atomic"
	"time"
	"unsafe"
)

// CSource 配置信息源
type CSource interface {
	// Get 获取指定key信息
	Get(serviceName, key string) (value string, err error)
}

// MConfig 动态配置handler,需要初始化里面的source来源
type MConfig struct {
	CSource
	watchList   []*mCWatchInfo // 需要定时watch的列表
	serviceName string         // 服务名
	ticker      int64          // 动态更新间隔时间,单位s
}

type MCWatchInfo struct {
	Def interface{}     // 自定义配置结构体,需要是指针类型,作为结构体和缺省值,不直接用来作为配置文件
	PT  *unsafe.Pointer // 地址信息,动态配置使用该地址进行转化
	Key string          // 监听的对应的key
}

// mCWatchInfo 需要定时watch的配置信息
type mCWatchInfo struct {
	def    interface{}                     // 自定义配置结构体,需要是指针类型,作为结构体和缺省值,不直接用来作为配置文件
	defMap map[string]*jsoniter.RawMessage // 将缺省值解析成 map[string]*jsoniter.RawMessage,使用默认值填充缺省值
	pt     *unsafe.Pointer                 // 地址信息,动态配置使用该地址进行转化
	key    string                          // 监听的对应的key
	curMD5 string                          // 当前key对应的值的md5,用来比较是否需要进行更新
}

// NewConfig new一个config配置的handler
func NewConfig(source CSource, watchList []*MCWatchInfo, opt ...DynamicOption) (*MConfig, error) {
	mc := &MConfig{
		CSource:     source,
		watchList:   make([]*mCWatchInfo, 0),
		serviceName: DefServiceName,
		ticker:      DefRefreshTicker,
	}

	// 添加其他配置文件
	for _, o := range opt {
		o(mc)
	}

	// 添加需要监控的信息
	err := mc.addWatchInfo(watchList...)
	if err != nil {
		log.Errorf("add watch info fail,err:%s", err.Error())
		return nil, err
	}

	return mc, nil
}

// Watch 对动态配置文件进行定时检查更新
func (c *MConfig) Watch() {
	// 后续定时更新配置文件
	go func() {
		for {
			time.Sleep(time.Duration(c.ticker) * time.Second)
			for i := range c.watchList {
				// 获取key对应的value
				nv, errT := c.Get(c.serviceName, c.watchList[i].key)
				if errT != nil {
					log.Errorf("%s get key %s value fail,err:%s", c.serviceName, c.watchList[i].key, errT.Error())
					continue
				}

				// 值为空的话直接不进行参数赋值
				if nv == "" {
					log.Debugf("%s get key %s value is nil", c.serviceName, c.watchList[i].key)
					continue
				}
				// 如果没有变化,则不用更新
				nMD5 := getMd5(nv)
				if c.watchList[i].curMD5 == nMD5 {
					continue
				}

				c.watchList[i].curMD5 = nMD5
				// 单独更新每一个key对应的配置文件
				safeWatch(c.watchList[i], nv)
			}
		}
	}()
	return
}

// safeWatch 处理每一个watch info
func safeWatch(w *mCWatchInfo, nv string) {
	// recover住panic
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("key:%s panic,%s", w.key, r)
		}
	}()

	// 把动态配置解析成 map[string]*jsoniter.RawMessage
	var tmp map[string]*jsoniter.RawMessage
	err := jsonCaseSensitive.Unmarshal([]byte(nv), &tmp)
	if nil != err {
		log.Errorf("unmarshal value fail,nv:%s", nv)
		return
	}

	// 补全缺省值
	if len(tmp) > 0 {
		for k, v := range w.defMap {
			if _, ok := tmp[k]; !ok {
				tmp[k] = v
			}
		}
	} else {
		tmp = w.defMap
	}

	// 再把map转bytes，可能有别的方法省略这一步
	jData, _ := jsonCaseSensitive.Marshal(&tmp)

	// 创建一个新的实例，避免读写冲突
	ni := reflect.New(reflect.ValueOf(w.def).Elem().Type())
	if err := jsonCaseSensitive.Unmarshal(jData, ni.Interface()); err != nil {
		log.Errorf("safeWatch key:%s json unmarshal failed:%v", w.key, err)
		return
	}

	log.Debugf("safeWatch key:%s config update:%+v", w.key, ni)
	atomic.StorePointer(w.pt, unsafe.Pointer(ni.Elem().UnsafeAddr()))
	return
}

// AddWatchInfo 给用户直接调用
func (c *MConfig) AddWatchInfo(info ...*MCWatchInfo) (err error) {
	err = c.addWatchInfo(info...)
	if err != nil {
		log.Errorf("add fail,err:%s", err.Error())
		return err
	}

	return
}

// addWatchInfo 添加需要watch的配置信息
func (c *MConfig) addWatchInfo(info ...*MCWatchInfo) (err error) {
	// 需要校验MCWatchList中的def为指针,否则直接报错
	for _, v := range info {
		if reflect.Ptr != reflect.TypeOf(v.Def).Kind() {
			err = fmt.Errorf("def must be a pointer")
			log.Errorf("WatchStruct fail:%s", err.Error())
			return err
		}

		// key不能为空
		if len(v.Key) <= 0 {
			err = fmt.Errorf("key can't be nil")
			log.Errorf("WatchStruct fail:%s", err.Error())
			return err
		}

		// 使用marshal和unmarshal获取默认值的RawMessage信息
		defMap := make(map[string]*jsoniter.RawMessage)
		bs, _ := jsonCaseSensitive.Marshal(v.Def)
		if err = jsonCaseSensitive.Unmarshal(bs, &defMap); err != nil {
			log.Errorf("get def map unmarshal fail,err:%s", err.Error())
			return err
		}

		tmp := &mCWatchInfo{
			def:    v.Def,
			defMap: defMap,
			pt:     v.PT,
			key:    v.Key,
		}
		// 第一次执行直接将默认值赋给pt
		up := unsafe.Pointer(reflect.ValueOf(tmp.def).Elem().UnsafeAddr())
		atomic.StorePointer(tmp.pt, up)
		// 第一次同步更新配置文件
		// 获取key对应的value
		nv, err := c.Get(c.serviceName, tmp.key)
		if err != nil {
			log.Errorf("%s get key %s value fail,err:%s", c.serviceName, tmp.key, err.Error())
			return err
		}
		// 单独更新每一个key对应的配置文件
		safeWatch(tmp, nv)

		c.watchList = append(c.watchList, tmp)
	}
	return
}

func getMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

const (
	DefRefreshTicker = 10    // 配置文件默认更新间隔时间 10s
	DefServiceName   = "mkg" // 默认服务名-作为value存储标志
)

var (
	// 大小写敏感的json序列化反序列化
	jsonCaseSensitive = jsoniter.Config{
		EscapeHTML:    true,
		CaseSensitive: true,
	}.Froze()
)

// DynamicOption 设置配置文件相关可选参数
type DynamicOption func(options *MConfig)

// AddServiceName 设置服务名
func AddServiceName(name string) DynamicOption {
	return func(DynamicOption *MConfig) {
		DynamicOption.serviceName = name
	}
}

// FreshTicker 设置配置动态刷新时间,单位秒
func FreshTicker(ticker int64) DynamicOption {
	return func(DynamicOption *MConfig) {
		DynamicOption.ticker = ticker
	}
}
