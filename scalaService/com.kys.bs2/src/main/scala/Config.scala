import scala.io.Source
import scala.collection.mutable
import java.io.FileNotFoundException
import java.io.IOException

object EnvConfiguration {
  def loadEnv(): Map[String, String] = {
    try {
      val envFile = Source.fromFile(".env")
      val lines = envFile.getLines().filter(_.nonEmpty)
      val envVars = lines.map { line =>
        val Array(key, value) = line.split("=", 2)
        key.trim -> value.trim
      }.toMap
      envFile.close()
      envVars
    } catch {
      case e: FileNotFoundException =>
        throw new RuntimeException("The .env file was not found. Please ensure it exists in the working directory.", e)
      case e: IOException =>
        throw new RuntimeException("Error reading the .env file.", e)
    }
  }

  val envVars: Map[String, String] = loadEnv()
  
  val dbUrl: String = envVars.getOrElse("DB_URL", "url")
  val dbUser: String = envVars.getOrElse("DB_USER", "user")
  val dbPassword: String = envVars.getOrElse("DB_PASSWORD", "password")
  val dbDriver: String = envVars.getOrElse("DB_DRIVER", "jdbc:mysql://localhost:3306/defaultdb")
}
