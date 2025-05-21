import org.zeromq.ZMQ;
import com.fasterxml.jackson.databind.ObjectMapper;
import java.util.*;

public class Coordinator {
    private final ZMQ.Context context;
    private final BullyElection election;
    private final List<String> servers;  
    private final ObjectMapper mapper = new ObjectMapper();

    public Coordinator(ZMQ.Context ctx, String id, List<String> servers) {
        this.context = ctx;
        this.servers = servers;
        Set<String> set = new HashSet<>(servers);
        election = new BullyElection(ctx, id, set);
    }

    public void run() throws Exception {
        election.startElection();
        if (!election.isCoordinator()) {
            System.out.println("[Coordinator] Not coordinator, exiting");
            return;
        }
        System.out.println("[Coordinator] Acting as coordinator, starting Berkeley sync");

        Map<String, ZMQ.Socket> reqSockets = new HashMap<>();
        for (String srv : servers) {
            ZMQ.Socket r = context.socket(ZMQ.REQ);
            r.connect("tcp://" + srv);
            reqSockets.put(srv, r);
        }

        while (true) {
            Map<String, Double> offsets = new HashMap<>();
            for (String srv : servers) {
                ZMQ.Socket r = reqSockets.get(srv);
                Map<String,Object> syncReq = Map.of(
                    "type", "SYNC_REQUEST",
                    "from_id", "coordinator",
                    "to_id", srv,
                    "payload", Collections.emptyMap()
                );
                byte[] reqBytes = mapper.writeValueAsBytes(syncReq);
                r.send(reqBytes, 0);
                byte[] repBytes = r.recv(0);
                Map<String,Object> repMsg = mapper.readValue(repBytes, Map.class);
                String sid = (String) repMsg.get("from_id");
                Double offset = ((Number)((Map)repMsg.get("payload")).get("offset")).doubleValue();
                offsets.put(sid, offset);
                System.out.println("[Berkeley] Reply from " + sid + ": offset=" + offset);
            }

            double avg = offsets.values().stream().mapToDouble(d -> d).average().orElse(0);
            for (String srv : servers) {
                ZMQ.Socket r = reqSockets.get(srv);
                double srvOffset = offsets.getOrDefault(srv, avg);
                Map<String,Object> adjust = Map.of(
                    "type", "SYNC_ADJUST",
                    "from_id", "coordinator",
                    "to_id", srv,
                    "payload", Map.of("adjust", avg - srvOffset)
                );
                byte[] adjBytes = mapper.writeValueAsBytes(adjust);
                r.send(adjBytes, 0);
                r.recv(0);  // ack
                System.out.println("[Berkeley] Sent SYNC_ADJUST to " + srv);
            }

            Thread.sleep(10000);
        }
    }

    public static void main(String[] args) throws Exception {
        if (args.length < 2) {
            System.err.println("Usage: Coordinator <myId> <host1:port1,host2:port2,...>");
            System.exit(1);
        }
        String myId = args[0];
        List<String> servers = List.of(args[1].split(","));
        ZMQ.Context ctx = ZMQ.context(1);
        new Coordinator(ctx, myId, servers).run();
    }
}
