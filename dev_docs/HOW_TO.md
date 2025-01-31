# How To

## How to add a new intake field

Adding a new intake field requires changes to both apm-data and apm-server repositories. 

### 1. Make changes in apm-data (this repo)

1. Add the new field to modeldecoder `model.go` such that the field is parsed from JSON.
   - Intake v2: [input/elasticapm/internal/modeldecoder/v2/model.go](/input/elasticapm/internal/modeldecoder/v2/model.go)
   - RUM v3: [input/elasticapm/internal/modeldecoder/rumv3/model.go](/input/elasticapm/internal/modeldecoder/rumv3/model.go)
2. Run `make generate` to generate the corresponding `model_generated.go` and JSON specification from the modified `model.go` in step 1.
3. Run `make update-licenses` to add license header to `model_generated.go` generated in step 2.
4. Add the new field to the corresponding file in `model/`, e.g. `span.go`.
   1. Add the field to the struct.
   2. Modify the `fields` function to include the field.
   3. Add a test to the corresponding `*_test.go` file.
5. Modify the modeldecoder `decoder.go` to map the modeldecoder model in step 1 to the internal model in step 4.
   - Intake v2: [input/elasticapm/internal/modeldecoder/v2/decoder.go](/input/elasticapm/internal/modeldecoder/v2/decoder.go)
   - RUM v3: [input/elasticapm/internal/modeldecoder/rumv3/decoder.go](/input/elasticapm/internal/modeldecoder/rumv3/decoder.go)
6. Run `make test`

### 2. Make changes in [apm-server](https://github.com/elastic/apm-server/)

1. Use the modified apm-data by replacing the apm-data dependency with the local one.
   - Run

         go mod edit -replace=github.com/elastic/apm-data=/path/to/your/apm-data
         make update
2. Modify [apmpackage](https://github.com/elastic/apm-server/tree/main/apmpackage) to add the field to Elasticsearch mapping.
   1. Find the corresponding data stream directory in `apmpackage/apm/data_stream/`. If the field applies to multiple data streams (e.g. field `service.name`), make sure all the corresponding data streams are updated. 
   2. Add the field to the YAML file in the data stream fields directory, e.g. `apmpackage/apm/data_stream/traces/fields/`.
        - Modify `ecs.yml`, if the field is defined in [ECS](https://www.elastic.co/guide/en/ecs/current/ecs-field-reference.html).
        - Otherwise, modify `fields.yml`.
   3. In case any changes to ingest pipelines and ILM policies are needed, they are inside the `apmpackage/apm/data_stream/<data_stream_name>/elasticsearch/` directory.
        - The common pipelines are defined in [apmpackage/cmd/genpackage/pipelines.go](https://github.com/elastic/apm-server/blob/main/apmpackage/cmd/genpackage/pipelines.go) and injected at package build (`make build-package`) time.
   4. Update apmpackage changelog `apmpackage/apm/changelog.yml`.

### 3. Test your changes with system test (in apm-server)

Modify apm-server system test to ensure the field works end-to-end.

1. Modify the input of the system test.
   - Intake v2: [`testdata/intake-v2/events.ndjson`](https://github.com/elastic/apm-server/blob/main/testdata/intake-v2/events.ndjson)
   - RUM v3: [`testdata/intake-v3/rum_events.ndjson`](https://github.com/elastic/apm-server/blob/main/testdata/intake-v3/rum_events.ndjson)
2. Run `make system-test` or only the specific tests.
    - Intake v2: [`systemtest/intake_test.go`](https://github.com/elastic/apm-server/blob/main/systemtest/intake_test.go)
    - RUM v3: [`systemtest/rum_test.go`](https://github.com/elastic/apm-server/blob/main/systemtest/rum_test.go)
3. System tests should fail as the received Elasticsearch documents do not match the expected documents because of the new field. If it doesn't fail, check the code.
4. Run `make update check-approvals` to review and accept the changes in the Elasticsearch documents.

### 4. Test your changes manually

See [apm-server TESTING.md](https://github.com/elastic/apm-server/blob/main/dev_docs/TESTING.md#manual-testing)

### 5. Finalize PRs

1. Create a PR in apm-data, and have it reviewed and merged.
2. In apm-server, bump apm-data dependency.
   - Run

         go mod edit -dropreplace=github.com/elastic/apm-data
         go get github.com/elastic/apm-data@main
         make update
3. Create a PR in apm-server.

### Example set of PRs to add an intake field:
- [apm-data PR](https://github.com/elastic/apm-data/pull/3)
- [apm-server PR](https://github.com/elastic/apm-server/pull/9850)

## How to map an OTel field

Mapping an OTel field is similar to adding a field to Intake.

### 1. Make changes in apm-data (this repo)

1. Modify OTel parsing code in [input/otlp](/input/otlp)
2. Add the new field to the corresponding file in `model/`, e.g. `span.go`.
   1. Add the field to the struct.
   2. Modify the `fields` function to include the field.
   3. Add a test to the corresponding `*_test.go` file.
3. Run `make test`

### 2. Make changes in [apm-server](https://github.com/elastic/apm-server/)

See the instructions in the same section under [How to add a new intake field](#how-to-add-a-new-intake-field)

### 3. Test your changes with system test (in apm-server)

Modify apm-server system test to ensure the field works end-to-end.

1. Modify the OTel system test to include the field if needed.
    - [`systemtest/otlp_test.go`](https://github.com/elastic/apm-server/blob/main/systemtest/otlp_test.go)
2. Run `make system-test` or only the above test.
3. System tests should fail as the received Elasticsearch documents do not match the expected documents because of the new field. If it doesn't fail, check the code.
4. Run `make update check-approvals` to review and accept the changes in the Elasticsearch documents.

### 4. Test your changes manually

See [apm-server TESTING.md](https://github.com/elastic/apm-server/blob/main/dev_docs/TESTING.md#manual-testing)

### 5. Finalize PRs
1. Same as [How to add a new intake field](#how-to-add-a-new-intake-field)

### Example PR:
- [apm-server PR](https://github.com/elastic/apm-server/pull/8334)
  - PR was done before apm-data repo was extracted from apm-server.
