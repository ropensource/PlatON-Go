package evm.lib;

import evm.beforetest.ContractPrepareTest;
import network.platon.autotest.junit.annotations.DataSource;
import network.platon.autotest.junit.enums.DataSourceType;
import network.platon.contracts.SafeMathMock;
import org.junit.Test;

import java.math.BigInteger;

public class SafeMathMockTest extends ContractPrepareTest {

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "max",
            author = "albedo", showName = "lib.SafeMathMockTest-最大值", sourcePrefix = "evm")
    public void testMax() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.load("0x20ec0ed750fd8fe9c1b7bbd63a1f21e37ddb42c9",web3j, transactionManager, provider);
            String contractAddress = using.getContractAddress();
            BigInteger result = using.max(new BigInteger("12"), new BigInteger("13")).send();
            collector.assertEqual(result, new BigInteger("13"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testMax failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "min",
            author = "albedo", showName = "lib.SafeMathMockTest-最小值", sourcePrefix = "evm")
    public void testMin() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SafeMathMock issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            BigInteger result = using.min(new BigInteger("12"), new BigInteger("13")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testMin failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "average",
            author = "albedo", showName = "lib.SafeMathMockTest-平均值", sourcePrefix = "evm")
    public void testAverage() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.load("0x5b9513089e960361162e6f8082eb136ae150fc83",web3j, transactionManager, provider);
            String contractAddress = using.getContractAddress();
            BigInteger result = using.average(new BigInteger("12"), new BigInteger("13")).send();
            collector.assertEqual(result, new BigInteger("12"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testAverage failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "add",
            author = "albedo", showName = "lib.SafeMathMockTest-无符号整型相加", sourcePrefix = "evm")
    public void testAdd() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.load("0xe63aa66f3c6e8a2c29d99193c0f5748466e662c0",web3j, transactionManager, provider);
            String contractAddress = using.getContractAddress();
            BigInteger result = using.add(new BigInteger("12"), new BigInteger("13")).send();
            collector.assertEqual(result, new BigInteger("25"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testAdd failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "mul",
            author = "albedo", showName = "lib.SafeMathMockTest-无符号整型相乘", sourcePrefix = "evm")
    public void testMul() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.load("0xce54ff31b3946ca8d548ad60656a6b91c5b94d36",web3j, transactionManager, provider);
            String contractAddress = using.getContractAddress();
            BigInteger result = using.mul(new BigInteger("12"), new BigInteger("13")).send();
            collector.assertEqual(result, new BigInteger("156"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testMul failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "sub",
            author = "albedo", showName = "lib.SafeMathMockTest-无符号整型相减", sourcePrefix = "evm")
    public void testSub() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.load("0x2a47e3309eb32530762c4b81e659613384f57993",web3j, transactionManager, provider);
            String contractAddress = using.getContractAddress();
            BigInteger result = using.sub(new BigInteger("12"), new BigInteger("13")).send();
            collector.assertEqual(result, new BigInteger("3963877391197344453575983046348115674221700746820753546331534351508065746944"), "checkout library function");
            result = using.sub(new BigInteger("13"), new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("1"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testSub failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }

    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "div",
            author = "albedo", showName = "lib.SafeMathMockTest-无符号整型相除", sourcePrefix = "evm")
    public void testDiv() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SafeMathMock issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            BigInteger result = using.div(new BigInteger("12"), new BigInteger("0")).send();
            collector.assertEqual(result, new BigInteger("3963877391197344453575983046348115674221700746820753546331534351508065746944"), "checkout library function");
            result = using.div(new BigInteger("13"), new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("1"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testDiv failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }
    @Test
    @DataSource(type = DataSourceType.EXCEL, file = "test.xls", sheetName = "mod",
            author = "albedo", showName = "lib.SafeMathMockTest-无符号整型除余", sourcePrefix = "evm")
    public void testMod() {
        try {
            prepare();
            SafeMathMock using = SafeMathMock.deploy(web3j, transactionManager, provider).send();
            String contractAddress = using.getContractAddress();
            String transactionHash = using.getTransactionReceipt().get().getTransactionHash();
            collector.logStepPass("SafeMathMock issued successfully.contractAddress:" + contractAddress + ", hash:" + transactionHash);
            BigInteger result = using.mod(new BigInteger("12"), new BigInteger("0")).send();
            collector.assertEqual(result, new BigInteger("3963877391197344453575983046348115674221700746820753546331534351508065746944"), "checkout library function");
            result = using.mod(new BigInteger("13"), new BigInteger("12")).send();
            collector.assertEqual(result, new BigInteger("1"), "checkout library function");
        } catch (Exception e) {
            collector.logStepFail("SafeMathMockTest testMod failure,exception msg:", e.getMessage());
            e.printStackTrace();
        }
    }
}
