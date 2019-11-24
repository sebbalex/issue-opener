package engines

import log "github.com/sirupsen/logrus"

// time="2019-11-18T01:05:04Z" level=error msg="Error parsing publiccode.yml for https://bitbucket.org/Comune_Venezia/iris/raw/master/publiccode.yml."
// time="2019-11-18T01:05:04Z" level=error msg="[Comune_Venezia/iris] invalid publiccode.yml: url: invalid repository URL: https://bitbucket.org/Comune_Venezia/iris/src/master\nlogo: HTTP GET returned 404 for https://bitbucket.org/Comune_Venezia/iris/raw/master/public/bann
// er.png; 200 was expected\ndescription/it/screenshots: HTTP GET returned 404 for https://bitbucket.org/Comune_Venezia/iris/raw/master/public/screenshot.png; 200 was expected"
// time="2019-11-18T01:05:04Z" level=error msg="Appending the bad file URL to the list: https://bitbucket.org/Comune_Venezia/iris/raw/master/publiccode.yml"
// time="2019-11-18T01:05:04Z" level=info msg="[RegioneUmbria/ecosistema-puppet] publiccode.yml found at https://raw.githubusercontent.com/RegioneUmbria/ecosistema-puppet/master/publiccode.yml"
// time="2019-11-18T01:05:04Z" level=error msg="Error parsing publiccode.yml for https://bitbucket.org/Comune_Venezia/whistleblowing/raw/master/publiccode.yml."
// time="2019-11-18T01:05:04Z" level=error msg="[Comune_Venezia/whistleblowing] invalid publiccode.yml: url: invalid repository URL: https://bitbucket.org/Comune_Venezia/whistleblowing/"

// TestBB func
func TestBB() {
	log.Debug("test bitbucket")
}
