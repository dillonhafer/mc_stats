import React, { Component } from 'react';
import './App.css';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import statMap from './statMap';

const DataRow = ({ label, value }) => {
  return (
    <tr>
      <td style={{ fontWeight: 'bold', color: '#777' }}>{label}</td>
      <td style={{ textAlign: 'right', fontWeight: 'bold' }}>
        {value || '😐'}
      </td>
    </tr>
  );
};

const Player = ({ player, data }) => {
  return (
    <div className="player">
      <table style={{ width: '100%', padding: '8px' }}>
        <tbody>
          <DataRow label="Died:" value={data['stat.deaths']} />
          <DataRow label="Chests Opened:" value={data['stat.chestOpened']} />
          <DataRow label="Animals Bred:" value={data['stat.animalsBred']} />
          <DataRow label="Left Game:" value={data['stat.leaveGame']} />
          <DataRow
            label="Villagers Talked To:"
            value={data['stat.talkedToVillager']}
          />
          {Object.entries(statMap).map(e => {
            const label = e[1];
            const key = e[0];
            return <DataRow key={key} label={label} value={data[key]} />;
          })}
        </tbody>
      </table>
    </div>
  );
};

class App extends Component {
  state = {
    loading: false,
    players: [],
    stats: [],
  };

  componentDidMount() {
    this.loadPlayers();
    this.loadStats();
    window.setInterval(this.loadStats, 5000);
  }

  loadStats = async () => {
    try {
      this.setState({ loading: true });

      const request = {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        method: 'GET',
      };

      const r = await fetch('/stats', request);
      const json = await r.json();
      if (json) {
        this.setState({ stats: json });
      }
    } catch (err) {
      console.warn(err);
    } finally {
      this.setState({ loading: false });
    }
  };

  loadPlayers = async () => {
    try {
      this.setState({ loading: true });

      const request = {
        headers: {
          Accept: 'application/json',
          'Content-Type': 'application/json',
        },
        method: 'GET',
      };

      const r = await fetch('/players', request);
      const json = await r.json();
      if (json) {
        this.setState({ players: json });
      }
    } catch (err) {
      console.warn(err);
    } finally {
      this.setState({ loading: false });
    }
  };

  players = () => {
    const { players, stats } = this.state;
    return players.map(p => {
      const stat = stats.find(s => {
        return s.UUID === p.uuid;
      }) || { data: {} };
      return (
        <TabPanel key={p.uuid}>
          <Player player={p} data={stat.data} />
        </TabPanel>
      );
    });
  };

  render() {
    const { loading } = this.state;
    const notLoading = !loading;
    return (
      <div className="App">
        {notLoading && this.players()}
        <Tabs defaultFocus={true}>
          <TabList>
            {this.state.players.map(({ name }) => <Tab key={name}>{name}</Tab>)}
          </TabList>
          {this.players()}
        </Tabs>
      </div>
    );
  }
}

export default App;
