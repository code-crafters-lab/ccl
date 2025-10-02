import { createConnectTransport, createGrpcWebTransport } from "@connectrpc/connect-web";
import { CategoryService } from "@cc-labs/proto/category/v1/category_service_pb";
import { createClient } from "@connectrpc/connect";
const transport = createConnectTransport({
  baseUrl: "http://localhost:8090/api/v1",
  useHttpGet: true,
});

// const transport = createGrpcWebTransport({
//   baseUrl: "http://localhost:8090",
// });

const categoryClient = createClient(CategoryService, transport);
const res = await categoryClient.listCategory({ pagination: { number: 1, size: 20 } });
console.log(res.total);
console.log(res.categories);
