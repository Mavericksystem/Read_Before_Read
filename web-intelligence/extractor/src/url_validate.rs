use std::net::{IpAddr, ToSocketAddrs}:

#[derive(Debug)]
pub enum ValidationError {
    BadScheme,
    UnresolvableHost,
    BlockedAddress(IpAddr),
    
}