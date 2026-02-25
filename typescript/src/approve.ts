import { rpc, signHash } from "./client";

const MAX_FEE = "1%";

async function main() {
  const res = await rpc("hl_buildApproveBuilderFee", { maxFeeRate: MAX_FEE });
  const sig = await signHash(res.result.hash);

  await rpc("hl_sendApproval", {
    nonce: res.result.nonce,
    signature: sig,
    maxFeeRate: MAX_FEE,
  });

  console.log("Builder fee approved.");
}

main();
