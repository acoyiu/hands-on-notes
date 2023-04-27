#!/bin/bash

# create the vschema
vtctldclient \
  ApplyVSchema \
  --server localhost:15999 \
  --vschema-file _db_vschema_alpha_initial.json \
  alpha

# wait for VSchema is ready
echo "wait 5s for VSchema"
sleep 5

# create the schema
vtctldclient \
  ApplySchema \
  --server localhost:15999 \
  --sql-file _db_create_alpha_schema.sql \
  alpha

echo "Schema is applied"