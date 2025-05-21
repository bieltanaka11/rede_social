import org.zeromq.ZMQ;
import java.util.Set;
import java.util.HashSet;

public class BullyElection {
    private final ZMQ.Context context;
    private final String id;              
    private final Set<String> peers;      
    private volatile boolean isCoordinator = false;

    public BullyElection(ZMQ.Context ctx, String id, Set<String> allServers) {
        this.context = ctx;
        this.id = id;
        this.peers = new HashSet<>();
        for (String peer : allServers) {
            if (peer.compareTo(id) > 0) {
                peers.add(peer);
            }
        }
    }

    public void startElection() {
        ZMQ.Socket sock = context.socket(ZMQ.PUB);
        sock.bind("tcp://*:" + (6000 + Integer.parseInt(id.replaceAll("\\D+", ""))));
        ZMQ.Socket sub = context.socket(ZMQ.SUB);
        for (String peer : peers) {
            sub.connect("tcp://" + peer + ":" + (6000 + Integer.parseInt(peer.replaceAll("\\D+", ""))));
        }
        sub.subscribe("ELECTION".getBytes(ZMQ.CHARSET));
        boolean higherResponded = false;


        String electionMsg = "ELECTION " + id;
        sock.send(electionMsg);
        System.out.println("[Bully] " + id + " sent ELECTION");

        long start = System.currentTimeMillis();
        while (System.currentTimeMillis() - start < 5000) { 
            byte[] reply = sub.recv(ZMQ.DONTWAIT);
            if (reply != null) {
                String msg = new String(reply, ZMQ.CHARSET);
                if (msg.startsWith("OK")) {
                    higherResponded = true;
                    break;
                }
            }
        }

        if (!higherResponded) {

            String coordMsg = "COORDINATOR " + id;
            sock.send(coordMsg);
            isCoordinator = true;
            System.out.println("[Bully] " + id + " is now coordinator");
        } else {
            System.out.println("[Bully] " + id + " yields, waiting coordinator");
        }

        sock.close();
        sub.close();
    }

    public boolean isCoordinator() {
        return isCoordinator;
    }
}
