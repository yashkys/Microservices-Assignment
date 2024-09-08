import java.time.LocalDateTime
import java.time.format.DateTimeFormatter

object Logger {
  
  // Define the formatter for date and time
  private val formatter = DateTimeFormatter.ofPattern("yyyy-MM-dd HH:mm:ss")

  // ANSI escape codes for colors
  private val Reset = "\u001B[0m"
  private val Red = "\u001B[31m"
  private val Green = "\u001B[32m"
  private val Yellow = "\u001B[33m"
  private val Blue = "\u001B[34m"

  // Function to get the current timestamp
  private def timestamp: String = LocalDateTime.now.format(formatter)

  // Log methods
  def info(message: String): Unit = {
    println(s"${Blue}INFO${Reset}: $timestamp - $message")
  }

  def success(message: String): Unit = {
    println(s"${Green}SUCCESS${Reset}: $timestamp - $message")
  }

  def warning(message: String): Unit = {
    println(s"${Yellow}WARNING${Reset}: $timestamp - $message")
  }

  def error(message: String): Unit = {
    println(s"${Red}ERROR${Reset}: $timestamp - $message")
  }
}