use std::net::{IpAddr, ToSocketAddrs}:

#[derive(Debug)]
pub enum ValidationError {
    BadScheme,
    UnresolvableHost,
    BlockedAddress(IpAddr),
}

impl std::fmt::Display for ValidationError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self{
            ValidationError::BadScheme => write!(f, "url must use http or https scheme"),
            ValidationError::UnresolvableHost => write!(f, "could not resolve host"),
            ValidationError::BlockedAddress(ip) => {
                write!(f, "resolved address {} is in a blocked range", ip)
            }
        }
    }
}