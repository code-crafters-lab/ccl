import { CategoryService } from "@cc-labs/proto/category/v1/category_service_pb";
import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";

export { createClientFor, toDate } from "./helpers.js";
export { NewAuthorizationBearerInterceptor } from "./interceptors.js";

// // TODO: Move this to `./protobuf.ts` and export it from there
export { create, fromJson, toJson } from "@bufbuild/protobuf";
export type { JsonObject } from "@bufbuild/protobuf";
export type { GenService } from "@bufbuild/protobuf/codegenv1";
export type { Duration, Timestamp } from "@bufbuild/protobuf/wkt";
export type { Client, Code, ConnectError } from "@connectrpc/connect";

export {
  TimestampSchema,
  timestampDate,
  timestampFromDate,
  timestampFromMs,
  timestampMs,
} from "@bufbuild/protobuf/wkt";

const transport = createConnectTransport({
  baseUrl: "http://localhost:8090",
});

const categoryClient = createClient(CategoryService, transport);

const res = await categoryClient.listCategory({ pagination: { number: 1, size: 20 } });
console.log(res.total);
console.log(res.categories);
