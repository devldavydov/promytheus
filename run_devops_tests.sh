set -e

echo "### Build"
cd cmd/agent
go build .
cd ../..

cd cmd/server
go build .
cd ../..

echo
echo "### Iteration 1"
./devopstest -test.v -test.run=^TestIteration1$ -agent-binary-path=cmd/agent/agent

echo
echo "### Iteration 2"
./devopstest -test.v -test.run=^TestIteration2[b]*$ -source-path=. -binary-path=cmd/server/server

echo
echo "### Iteration 3"
