use hyperliquid_api_examples::Client;
use serde_json::json;

const COIN: &str = "BTC";

#[tokio::main]
async fn main() {
    let client = Client::from_env();

    let mid = client.get_mid(COIN).await;
    if mid == 0.0 {
        eprintln!("Could not fetch {COIN} mid price");
        std::process::exit(1);
    }

    let sz = format!("{:.5}", 11.0 / mid);
    let rest_px = format!("{}", (mid * 0.97) as u64);

    println!("{COIN} mid: ${mid:.2}");
    println!("Placing resting BUY {sz} @ {rest_px} (GTC, 3% below mid)\n");

    let res = client
        .exchange(&json!({
            "action": {
                "type": "order",
                "orders": [{"asset": COIN, "side": "buy", "price": rest_px, "size": sz, "tif": "gtc"}],
            },
        }))
        .await;

    let hash = res["hash"].as_str().unwrap();
    let sig = client.sign_hash(hash).await;

    let result = client
        .exchange(&json!({
            "action": res["action"],
            "nonce": res["nonce"],
            "signature": sig,
        }))
        .await;

    let exchange_resp = &result["exchangeResponse"];
    let statuses = exchange_resp["response"]["data"]["statuses"]
        .as_array()
        .expect("No statuses in response");

    let oid = statuses
        .iter()
        .find_map(|s| s["resting"]["oid"].as_u64())
        .expect("Could not extract OID from resting order");

    println!("Order resting (OID: {oid})");
    println!("Cancelling...\n");

    let cancel_action = json!({
        "type": "cancel",
        "cancels": [{"a": COIN, "o": oid}],
    });

    let res = client.exchange(&json!({"action": cancel_action})).await;

    let hash = res["hash"].as_str().unwrap();
    let sig = client.sign_hash(hash).await;

    let cancel_result = client
        .exchange(&json!({
            "action": cancel_action,
            "nonce": res["nonce"],
            "signature": sig,
        }))
        .await;

    println!(
        "{}",
        serde_json::to_string_pretty(&cancel_result["exchangeResponse"]).unwrap()
    );
    println!("\nOrder cancelled.");
}
