test: test-units test-acceptance

test-units:
	echo "running unit tests"
	ginkgo .
	ginkgo -r controller/
	ginkgo -r reconciler/

test-acceptance:
	echo "running acceptance tests"
	ginkgo -r acceptance_test/
