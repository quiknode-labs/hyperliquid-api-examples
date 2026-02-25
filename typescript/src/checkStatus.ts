import { rpc, wallet } from "./client";

async function main() {
  const res = await rpc("hl_getMaxBuilderFee", { user: wallet.address });
  console.log(JSON.stringify(res.result, null, 2));
}

main();
