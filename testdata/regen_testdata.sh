#!/bin/bash

prefixes="basic codeblock ignore repeat special"
for prefix in $prefixes; do
    mdtoc -out=testdata/${prefix}_toc.md testdata/${prefix}.md
done

# handle separately with custom heading flag
mdtoc -with-toc-heading -toc-heading="Contents" -out=testdata/custom_heading_toc.md testdata/custom_heading.md
