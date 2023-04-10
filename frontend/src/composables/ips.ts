export const isValidIpv4 = (ipaddress: string) => {
  const ipv4Pattern =
    /^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
  return ipv4Pattern.test(ipaddress);
};

export const isValidIpv6 = (ipaddress: string) => {
  const ipv6Pattern = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/;
  return ipv6Pattern.test(ipaddress);
};

export const isValidIp = (ipaddress: string) =>
  isValidIpv4(ipaddress) || isValidIpv6(ipaddress);
