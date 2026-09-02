use std::net::{IpAddr, ToSocketAddrs}:

#[derive(Debug)]
pub enum ValidationError {
    BadScheme,
    UnresolvableHost,
    BlockedAddress(IpAddr),
}

impl std::fmt::Display for ValidatioError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self{
            
        }
    }
}