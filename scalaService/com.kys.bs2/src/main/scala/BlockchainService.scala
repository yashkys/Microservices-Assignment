import com.kys.bs2.transaction.transaction.TransactionServiceGrpc.TransactionService
import scala.concurrent.Future
import com.kys.bs2.transaction.transaction.SubmitResponse
import com.kys.bs2.transaction.transaction.TransactionReceipt
import java.sql.PreparedStatement
import java.sql.Connection
import java.sql.DriverManager
import java.sql.Statement

class BlockchainServiceImpl extends TransactionService {
    // Database connection parameters
    val url = EnvConfiguration.dbUrl
    val user = EnvConfiguration.dbUser
    val password = EnvConfiguration.dbPassword
    val driver = EnvConfiguration.dbDriver
    var connection: Connection = null
    var statement: Statement = null

    // Initialize database connection
    def initDBConnection(): Unit = {
      try{
        Class.forName(driver)
        connection = DriverManager.getConnection(url, user, password)
        // Start transaction
        connection.setAutoCommit(false)
        /* Note not required in production Onlly use SQL quesries in the db directly */
        statement = connection.createStatement()
        val createBlockTransactionsTable = """
        CREATE TABLE IF NOT EXISTS block_transactions (
            block_number VARCHAR(255) PRIMARY KEY,
            receipt_data TEXT,
            transaction_count INT
        )
        """
        statement.executeUpdate(createBlockTransactionsTable)
        statement.close()
        Logger.info("Database connected")
      } catch {
              case e: Exception =>
                Logger.error("Database access denied")
          }

    }

    // Method to insert or update transaction receipt
    def insertOrUpdateTransactionReceipt(receipt: TransactionReceipt): Unit = {    
      val blockNumber = try {
        Integer.parseInt(receipt.blockNumber.substring(2), 16) // Convert hex string to integer
      } catch {
        case e: NumberFormatException =>
          Logger.error(s"Failed to parse blockNumber: ${e.getMessage}")
          0 // or handle as needed
      }
      try {
        val receiptData = receipt.toString  // Serialize to string or JSON
        val transactionCount = receipt.logs.size

        // SQL to insert or update the transaction receipt for a block
        val sql = """
          INSERT INTO block_transactions (block_number, receipt_data, transaction_count)
          VALUES (?, ?, ?)
          ON DUPLICATE KEY UPDATE
            receipt_data = CONCAT(receipt_data, '\n', VALUES(receipt_data)),
            transaction_count = transaction_count + VALUES(transaction_count)
        """

        val statement: PreparedStatement = connection.prepareStatement(sql)
        statement.setInt(1, blockNumber)
        statement.setString(2, receiptData)
        statement.setInt(3, transactionCount)

        statement.executeUpdate()
        Logger.success("Inserted data in db")

        // Commit transaction
        connection.commit()
      } catch {
          case e: Exception =>
            // Rollback in case of error
            connection.rollback()
            Logger.error(s"Failed to insert or update transaction receipt: ${e.getMessage}")
        } finally {
          connection.setAutoCommit(true)
          connection.close()
        }
      /*catch {
          case e: Exception =>
            Logger.error(s"Failed to insert data in db: ${e.getMessage}")
      }*/
    }

    override def submitTransactionReceipt(request: TransactionReceipt): Future[SubmitResponse] = {
        try {
            if (connection == null || connection.isClosed) initDBConnection()
            insertOrUpdateTransactionReceipt(request)
            val response = SubmitResponse(success = true)
            Future.successful(response)
        } catch {
            case e: Exception =>
                Logger.error(s"Unexpected: ${e.getMessage()}")
                val response = SubmitResponse(success = false)
                Future.successful(response)
        }
  }
}