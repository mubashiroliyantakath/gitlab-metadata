# gitlab-metadata

`gitlab-metadata` is a Go based tool for extracting metadata in GitLab CI/CD pipelines, especially prior to any Docker image build jobs or jobs that require build versions. This project is inspired by [Docker Metadata Action](https://github.com/marketplace/actions/docker-metadata-action) which provides a universal and easy to use tool for generating relevant metadata for building Docker images.

This expected to be run only on Gitlab CI pipelines as it depends on environment variables like
- `CI_PIPELINE_IID`
- `CI_COMMIT_SHORT_SHA`
- `CI_COMMIT_REF_SLUG`
- `CI_COMMIT_BRANCH`
- `CI_PROJECT_DIR`

being present during its execution.


## Customizing

Inputs for `metadata` is supplied via the environment variables`MA_INPUT`

| Name       | Type | Description                                              |
| ---------- | ---- | -------------------------------------------------------- |
| `MA_INPUT` | List | List of [inputs](#ma_input) as key-value pair CSV format |


### `MA_INPUT`

`MA_INPUT` specifies the kind of metadata that are to be generated. This is in a key-value pair CSV format to make it compatible with environment variables.

```yaml
MA_INPUT: |
    type=semver
    type=ref,event=branch
    type=sha
```

Each entry is defined by a `type`, which are

- [`type=semver`](#typesemver)
- `type=ref`
- `type=sha`

(if no input is supplied all the above types are generated)

and global attributes
- `enabled=<true|false>` Enable/Disables generation of the corresponding types (default=true)
- `suffix=<string>` Add suffix to the generated version or string (default `""`)
- `prefix=<string>` Add prefix to the generated version or string (default `""`)

An example with global attributes:
```yaml
MA_INPUT: |
    type=semver
    type=sha,suffix=,prefix=sha-,enabled=true
```

## Output
The output is a file named `metadata-out.env` which is compatible with the Gitlab CI [artifacts:reports:dotenv](https://docs.gitlab.com/ee/ci/yaml/artifacts_reports.html#artifactsreportsdotenv) artifacts (the key value pairs in these artifacts will be injected into the environment of the subsequent jobs automatically)

An example of the output generated:

```env
VERSION=1.1.2
TAGS=1.1.2,metadata-test,ff89123
```
where

- `1.1.2` is semver type generated
- `metadata-test` is the `type=ref,event=branch` type generated
- `ff89123` is the `type=sha` generated

The VERSION variable will be safe to build any projects that accepts the [standard semver](https://semver.org/) versions during the build process (like dotnet) and is generated only if the `type=semver` is enabled.

> This doesn't suppport including build metadata yet.

With `type=semver` enabled, it either automatically processes the commit tag being pushed or if it is not a tag pipeline, creates a transient version which will be a combination of the latest semver standard tag and the `CI_PIPELINE_IID`.


## Inputs

### `type=semver`
Including `type=semver` will generate `VERSION` and the corresponding tag in the `TAGS` section of the output. It either transforms the commit tag of the pipeline or automatically generates a version number by fetching the latest tag and appending the `CI_PIPELINE_IID` to it like `<LATEST_RELEASE_SEMVER_TAG>.<CI_PIPELINE_IID>`

The `VERSION` should be safe to use when building C# projects. And the items in the tags variable can be used to tag docker images being built for the project.

| Git tag           | VERSION Output  |
| ----------------- | --------------- |
| `v1.2.3`          | `1.2.3`         |
| `Release_1_2_3`   | `1.2.3`         |
| `v1.2.4.alpha07`  | `1.2.4-alpha07` |
| `v1_2_4-alpha_07` | `1.2.4-alpha07` |

Automatically generated semver output in case of a non-tag pipeline:
`1.2.3.113`
> This behaviour might change to create versions that are more aligned with the [standard semver](https://semver.org/) grammar by including build metadata information.

### `type=ref,event=branch`

This tag is only generated on pipelines that run for branches. This will be a slug of the corresponding branch.

| Branch       | Tag          |
| ------------ | ------------ |
| `feature/ci` | `feature-ci` |

### `type=sha`

This tag can be generated for all pipelines and this will be the short sha for the commit obtained from the `CI_COMMIT_SHORT_SHA` variable provided by the Gitlab CI.

Example, `bf411347`


## Usage

The following job definition could be used to propagate VERSION and TAGS variable for all subsequent jobs.

```yaml
generate-metadata:
  stage: metadata
  image: mubashiro/gitlab-metadata:0.1.2
  script:
    - /app/gitlab-metadata
  artifacts:
    expire_in: 1 day
    reports:
      dotenv: metadata-out.env
```
