use salvo::compression::Compression;
use salvo::prelude::*;

#[tokio::main]
async fn main() {
    tracing_subscriber::fmt().init();

    let base_dir = std::env::current_exe()
        .unwrap()
        .join("../../../static")
        .canonicalize()
        .unwrap();
    println!("Base Dir: {:?}", base_dir);

    let router = Router::new()
       .push(
            Router::with_hoop(Compression::new().enable_zstd(CompressionLevel::Minsize))
                .path("")
                .get(StaticFile::new(base_dir.join("data.json"))),
        );

    //没起效果还是 gzip encoding 在response headers里,是个bug需要fork
    
    let acceptor = TcpListener::new("127.0.0.1:5800").bind().await;
    Server::new(acceptor).serve(router).await;
}