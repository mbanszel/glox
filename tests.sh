#!/bin/bash
set -e

for a_test in examples/test_*.glox; do
  ./glox $a_test
done

