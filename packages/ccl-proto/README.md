# CCL Proto

This package provides the Protocol Buffers (proto) definitions used by ccl projects. It includes the proto files and generated code for interacting with ccl's gRPC APIs.

## Installation

To install the package, use pnpm or yarn:

```sh
pnpm install @cc-labs/proto
```

or

```sh
yarn add @zitadel/proto
```

## Usage

To use the proto definitions in your project, import the generated code:

```ts
import { Organization } from "@zitadel/proto/zitadel/org/v2/org_pb";

const org: Organization | null = await getDefaultOrg();
```

