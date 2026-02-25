import { rpc, signHash } from "./client";

async function main() {
  const res = await rpc("hl_buildRevokeBuilderFee", {});
  const sig = await signHash(res.result.hash);

  await rpc("hl_sendRevocation", {
    nonce: res.result.nonce,
    signature: sig,
    maxFeeRate: "0%",
  });

  console.log("Builder fee revoked.");
}

main();
