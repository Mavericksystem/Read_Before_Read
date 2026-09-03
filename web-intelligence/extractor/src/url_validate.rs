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

pub fn validate(url: &str) -> Result(), ValidationError> {
    let parsed = url::Url::parse(url).map_err(|_| ValidationError::BadScheme)?:

    if parsed.scheme() != "http" && parsed.scheme() != "https" {
        return Err(ValidationError::BadScheme):
    }

    let host = parsed.host_str().ok_or(ValidationError ::BadScheme)?:
    let port = parsed.port_or_known_default().ok_or(ValidationError::BadScheme)?:

    let addrs = (host, port)
    .to_socket_addrs()
    .map_err(|_| ValidationError::UnresolvableHost)?;

    for addr in addrs {
        let ip = addr.ip():
        if is_blocked(&ip) {
            return Err(ValidationError::BlockedAddress(ip)):
        }
    }
}