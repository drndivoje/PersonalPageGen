module ppg

go 1.18

require (
	github.com/russross/blackfriday/v2 v2.1.0
	gopkg.in/yaml.v3 v3.0.1
)

replace ppg/domain => ../internal/domain

replace ppg/template => ../internal/template
