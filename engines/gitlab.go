package engines

import log "github.com/sirupsen/logrus"

// time="2019-11-18T01:05:05Z" level=error msg="Error parsing publiccode.yml for https://gitlab.com/fusslab/fuss/raw/master/publiccode.yml."
// time="2019-11-18T01:05:05Z" level=error msg="[fusslab/fuss] invalid publiccode.yml: url: invalid repository URL: https://work.fuss.bz.it/projects"
// time="2019-11-18T01:05:05Z" level=error msg="Appending the bad file URL to the list: https://gitlab.com/fusslab/fuss/raw/master/publiccode.yml"

// TestGL func
func TestGL() {
	log.Debug("test gitlab")
}
