package traffic_http_fast

////Certs 证书集合结构体
//type Certs struct {
//	certs map[string]*tls.Certificate
//}
//
////Get 获取证书
//func (c *Certs) Get(hostName string) (*tls.Certificate, bool) {
//	if c == nil || len(c.certs) == 0 {
//		return nil, true
//	}
//	cert, has := c.certs[hostName]
//	if has {
//		return cert, true
//	}
//	hs := strings.Split(hostName, ".")
//	if len(hs) < 1 {
//		return nil, false
//	}
//
//	cert, has = c.certs[fmt.Sprintf("*.%s", strings.Join(hs[1:], "."))]
//	return cert, has
//}
//
//func newCerts(data map[string]*tls.Certificate) *Certs {
//	return &Certs{certs: data}
//}
