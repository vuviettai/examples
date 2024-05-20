use grpc_server::{proto::route_guide_server::RouteGuideServer, server::RouteService};
use tonic::transport::Server;
use tracing::info;

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    // construct a subscriber that prints formatted traces to stdout
    let subscriber = tracing_subscriber::FmtSubscriber::new();
    // use that subscriber to process traces emitted after this point
    tracing::subscriber::set_global_default(subscriber)?;
    let addr = "0.0.0.0:50051".parse()?;
    info!("Start grpc at address {:?}", &addr);
    let service = RouteService::default();

    Server::builder()
        .add_service(RouteGuideServer::new(service))
        .serve(addr)
        .await?;
    Ok(())
}
