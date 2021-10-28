import java.util.Scanner;

public class main {

    public static void main (String[] args) {

        System.out.println("Welcome to Big Data Processing ToolBox");
        System.out.println("Author: Zhengyu HU  Andrew ID: zhengyuh  Course: 14848");
        System.out.println("Please select the application to run");
        System.out.println("1. Jupyter Notebook");
        System.out.println("2. Apache Hadoop");
        System.out.println("3. Apache Spark");
        System.out.println("4. SonarQube & SonarScanner");
        System.out.println("Type the index of application (no-index input to stop)");
        String originalIP = "";

        Scanner scan = new Scanner(System.in);
        boolean isQuit = false;
        while (scan.hasNext()) {
            String str = scan.next();
            switch(str) {
                case "1":
                    String ip1 = originalIP + "";
                    System.out.println("Jupyter Notebook URL is " + ip1);
                    break;
                case "2":
                    String ip2 = originalIP + "";
                    System.out.println("Apache Hadoop URL is " + ip2);
                    break;
                case "3":
                    String ip3 = originalIP + "";
                    System.out.println("Apache Spark URL is " + ip3);
                    break;
                case "4":
                    String ip4 = originalIP + "";
                    System.out.println("SonarQube & SonarScanner URL is " + ip4);
                    break;
                default:
                    System.out.println("Quit the ToolBox");
                    isQuit = true;
                    break;
            }
            if (isQuit)
                break;
        }
        scan.close();
    }
}
