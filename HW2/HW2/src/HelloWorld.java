public class HelloWorld {
    public static void main(String[] args) {

        int count = 0;
        try {
            while (true) {
                Thread.sleep(2 * 1000);
                System.out.println("Hello World from Docker! (zhengyuh)");
            }
        } catch (InterruptedException e) {
            e.printStackTrace();
        }

    }
}
