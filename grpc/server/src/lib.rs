// pub mod proto {
//     mod routeguide;

use std::pin::Pin;

use proto::RouteNote;
use tokio_stream::Stream;
use tonic::Status;

//     //include!("./proto/routeguide.rs");
//     pub use route_guide_client::RouteGuideClient;
//     pub use route_guide_server::{RouteGuide, RouteGuideServer};
// }
pub mod proto;
// pub use proto::route_guide_client::RouteGuideClient;
// pub use proto::route_guide_server::{RouteGuide, RouteGuideServer};
pub mod server;

pub type ResponseStream = Pin<Box<dyn Stream<Item = Result<RouteNote, Status>> + Send>>;
