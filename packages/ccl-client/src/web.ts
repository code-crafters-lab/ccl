import { createGrpcWebTransport, type GrpcWebTransportOptions } from "@connectrpc/connect-web";
import { NewAuthorizationBearerInterceptor } from "./interceptors.js";

/**
 * Create a client transport using grpc web with the given token and configuration options.
 * @param token
 * @param opts
 */
export function createClientTransport(token: string, opts: GrpcWebTransportOptions) {
  return createGrpcWebTransport({
    ...opts,
    interceptors: [...(opts.interceptors || []), NewAuthorizationBearerInterceptor(token)],
  });
}
