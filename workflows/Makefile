all: \
	aws-clusters.alfredworkflow \
	aws-services.alfredworkflow \
	aws-taskdefs.alfredworkflow \
	circle-ci.alfredworkflow \
	github.alfredworkflow \
	datadog.alfredworkflow \
	yandex-translate.alfredworkflow

aws-clusters.alfredworkflow: aws-clusters/*
	cd aws-clusters; zip ../$@ *

aws-services.alfredworkflow: aws-services/*
	cd aws-services; zip ../$@ *

aws-taskdefs.alfredworkflow: aws-taskdefs/*
	cd aws-taskdefs; zip ../$@ *

circle-ci.alfredworkflow: circle-ci/*
	cd circle-ci; zip ../$@ *

github.alfredworkflow: github/*
	cd github; zip ../$@ *

datadog.alfredworkflow: datadog/*
	cd datadog; zip ../$@ *

yandex-translate.alfredworkflow: yandex-translate/*
	cd yandex-translate; zip ../$@ *
