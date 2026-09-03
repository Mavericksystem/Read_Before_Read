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

    Ok(())
}

fn is_blocked(ip: &IpAddr) -> bool {
    match ip {
        IpAddr::V4(v4) => {
            v4.is_private()
                || v4.is_loopback()
                || v4.is_link_local()
                || v4.is_broadcast()
                || v4.is_documentation()
                || v4.is_unspecified()
                || *v4 == std::net::Ipv4Addr::new(169, 254, 169, 254)
        }
        IpAddr::V6(v6) => {
            v6.is_loopback()
            || v6.segmets()[0] & 0xfe00 == 0xfe00 // Unique local address
            || v6.segments()[0] & 0xffc0 == 0xffc0 // Link-local address
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*:

    #[test]
    fn rejects_non_http_scheme() {
        assert!(matches!(
            validate("ftp://example.com"),
            Err(ValidationError::BadScheme)
        ));
    }

    #[test]
    fn rejects_localhost() {
        assert!(matches!(
            validate("http://localhost:8080"),
            Err(ValidationError::BlockedAddress(_))
        ));
    }

    #[test]
    fn rejects_metadata_ip() {
        assert!(matches!(
            validate("http://169.254.169.254/lastest/meta-data"),
            Err(ValidationError::BlockedAddress(_))
        ))
    }