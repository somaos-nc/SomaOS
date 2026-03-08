#!/bin/bash
# TEST/04_cljv_test.sh
# Dedicated Verification for ClojureV Linter & Multi-Target Transpiler

echo "--- [$1/13] RUNNING CLOJUREV TESTS (Linter & Transpiler) ---"
cd soma_os_go && go test -v ./pkg/clojurev/...
EXIT_CODE=$?
cd ..

if [ $EXIT_CODE -ne 0 ]; then
    echo "ClojureV Tests failed."
    exit $EXIT_CODE
fi
echo "ClojureV Tests passed."
