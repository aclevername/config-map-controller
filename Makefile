test: test-units test-acceptance

test-units:
	echo "running unit tests"
	ginkgo -r controller/

test-acceptance:
	echo "running acceptance tests"
	ginkgo -r acceptance_test/
