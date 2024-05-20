use crate::{
    proto::{
        route_guide_server::{RouteGuide, RouteGuideServer},
        RouteNote,
    },
    ResponseStream,
};
use tokio::sync::mpsc;
use tokio_stream::{wrappers::UnboundedReceiverStream, StreamExt};
use tonic::{transport::Server, Response};
use tracing::{error, info, span};
// pub mod message {
//     tonic::include_proto!("message");
// }

#[derive(Default, Clone)]
pub struct RouteService {}

#[tonic::async_trait]
impl RouteGuide for RouteService {
    type RouteChatStream = ResponseStream;
    /*
     * Consensus client init a duplex streaming connection to send external transaction
     * and to receives consensus output.
     * External trasaction contains a namespace field and a content in byte array
     */
    async fn route_chat(
        &self,
        request: tonic::Request<tonic::Streaming<RouteNote>>,
    ) -> std::result::Result<tonic::Response<Self::RouteChatStream>, tonic::Status> {
        info!("RouteService::route_chat");
        let mut in_stream = request.into_inner();
        let (tx_consensus, rx_consensus) = mpsc::unbounded_channel();
        // let service = self.clone();
        let _handle = tokio::spawn(async move {
            //let service = consensus_service;
            while let Some(client_message) = in_stream.next().await {
                match client_message {
                    Ok(route_note) => {
                        info!("Received route note {:?}", &route_note);
                        tx_consensus.send(Ok(route_note));
                        // let _handle_res =
                        //     service.handle_consensus_transaction(transaction_in).await;
                    }
                    Err(err) => {
                        error!("{:?}", err);
                    }
                }
            }
        });
        let out_stream = UnboundedReceiverStream::new(rx_consensus);

        Ok(Response::new(Box::pin(out_stream) as Self::RouteChatStream))
    }
}

#[tokio::main]
pub async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse()?;
    let service = RouteService::default();

    Server::builder()
        .add_service(RouteGuideServer::new(service))
        .serve(addr)
        .await?;

    Ok(())
}
