test: test-units test-acceptance

test-units:
	echo "running unit tests"
	ginkgo .
	ginkgo -r controller/
	ginkgo -r processor/

test-acceptance:
	echo "running acceptance tests"
	ginkgo -r acceptance_test/
