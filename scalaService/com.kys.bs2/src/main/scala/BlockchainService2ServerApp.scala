import io.grpc.ServerBuilder
import com.kys.bs2.transaction.transaction.TransactionServiceGrpc

object BlockchainService2Server {
    import scala.concurrent.ExecutionContext
    def main(args: Array[String]): Unit = {
        try {
            val port = 50051
            val server = ServerBuilder
            .forPort(port)
            .addService(TransactionServiceGrpc.bindService(new BlockchainServiceImpl, ExecutionContext.global))
            .build()
            .start()
            Logger.info("Server is running on port 50051")
            
            server.awaitTermination()
        } catch {
            case e: Exception =>
                var  msg = e.getMessage()
                if(msg == null){
                Logger.error("Server stopped")
                } else {
                Logger.error(s"Failed to start server.\n ${e.getMessage}")
                }
        }
    }
}