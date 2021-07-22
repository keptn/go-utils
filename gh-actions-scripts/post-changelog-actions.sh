#!/bin/bash

echo "Removing sign-off messages from changelog..."
for file in ['CHANGELOG.md' 'RELEASE-BODY.md']; do
  if [ -f "$file" ]; then
    sed ':a;N;$!ba;s/\nSigned-off-by: [A-Za-z0-9 ]*<.*@.*>\n//' "$file"
  fi
done
