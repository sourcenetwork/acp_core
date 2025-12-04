#!/usr/bin sh

buf generate
chown -R $HOST_USER ../pkg/types ../pkg/errors ../ui/src/types/proto-js 2>/dev/null || true
