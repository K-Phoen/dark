module github.com/K-Phoen/dark

go 1.13

require (
	github.com/K-Phoen/grabana v0.4.1
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.0.0-20200319202348-eb909d5fe0e7
	k8s.io/apimachinery v0.0.0-20200319202151-147abd67b880
	k8s.io/client-go v0.0.0-20200319202630-365234d2fcf0
	k8s.io/code-generator v0.0.0-20200319201949-6bb2b634cece
	k8s.io/klog v1.0.0
)

replace github.com/K-Phoen/grabana => ../grabana
