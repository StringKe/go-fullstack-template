import { createClient } from "@connectrpc/connect";
import { createConnectTransport } from "@connectrpc/connect-web";
import { TestService } from "./gen/v1/test_service_pb.ts";

const transport = createConnectTransport({
  baseUrl: "http://127.0.0.1:21421/rpc",
});

const client = createClient(TestService, transport);

function App() {
  const test1 = async () => {
    client
      .test1({})
      .then((res) => {
        console.log(res.message);
      })
      .catch((err) => {
        console.log(err);
      });
  };

  const test2 = async () => {
    client
      .test2({ name: "test" })
      .then((res) => {
        console.log(res.message);
      })
      .catch((err) => {
        console.log(err);
      });
  };

  const test3 = async () => {
    for await (const res of client.test3({ name: "test" })) {
      console.log(res.message);
    }
  };

  return (
    <>
      <button onClick={test1}>click to test1 (just send a request)</button>
      <button onClick={test2}>
        click to test2 (send a request with a name)
      </button>
      <button onClick={test3}>click to test3 (stream)</button>
    </>
  );
}

export default App;
